package lists

import (
	"math/rand"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/rwlist/youtube/internal/global"
	"github.com/rwlist/youtube/internal/models"
	"github.com/rwlist/youtube/internal/proto"
	"github.com/rwlist/youtube/internal/repos"
	"github.com/rwlist/youtube/internal/ytsync"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type testContext struct {
	db           *gorm.DB
	usersRepo    *repos.Users
	catalogsRepo *repos.CatalogLists
	globalDir    *global.Directory
	catalog      *Catalog
	user         *models.User
	liked        *LikedEngine
	likesCursor  *ytsync.MockCursor
}

func setupDatabase(t *testing.T) (x *testContext, cleanup func()) {
	connstr := os.Getenv("TEST_CONNSTR")
	if connstr == "" {
		t.Skip("TEST_CONNSTR not set")
	}

	migrator := newDatabaseMigrator(t, connstr)
	migrator.makeFreshDatabase()

	db := migrator.connectNew()

	err := models.AutoMigrate(db)
	if err != nil {
		t.Fatal("failed to migrate models", err)
	}

	usersRepo := repos.NewUsers(db)
	catalogsRepo := repos.NewCatalogLists(db)

	likesCursor := &ytsync.MockCursor{
		Rand: rand.New(rand.NewSource(42)),
	}

	globalDir := &global.Directory{
		LikedSync:    ytsync.NewLikedSync(likesCursor),
		CatalogsRepo: catalogsRepo,
	}

	user := &models.User{}
	assert.NoError(t, usersRepo.Save(user))

	catalog := NewCatalog(catalogsRepo, db, globalDir)
	assert.NoError(t, catalog.UserHook(user))

	likedEngine, err := catalog.GetList(user.ID, "liked")
	assert.NoError(t, err)

	x = &testContext{
		db:           db,
		usersRepo:    usersRepo,
		catalogsRepo: catalogsRepo,
		globalDir:    globalDir,
		catalog:      catalog,
		user:         user,
		liked:        likedEngine.(*LikedEngine),
		likesCursor:  likesCursor,
	}

	return x, func() {
		defer migrator.close()

		db, err := db.DB()
		if err != nil {
			t.Fatal("failed to get sql db", err)
		}
		err = db.Close()
		if err != nil {
			t.Fatal("failed to close sql db", err)
		}
	}
}

func dumpLists(t *testing.T, x *testContext) {
	lists, err := x.catalog.UserLists(x.user.ID)
	assert.NoError(t, err)

	for _, list := range lists {
		info, err := list.Info()
		assert.NoError(t, err)

		t.Log(spew.Sdump(info))

		if dbg, ok := list.(debugger); ok {
			t.Log(spew.Sdump(dbg.debugInfo()))
		}

		t.Log(spew.Sdump(list.ListItems()))
	}
}

//nolint:unparam
func fetchVerify(t *testing.T, liked *LikedEngine, youtubeItems []models.YoutubeData) (arr []models.LikedModel, xords []string) {
	items, err := liked.ListItemsTyped()
	assert.NoError(t, err)

	// verify that items are the same, in the same order
	assert.Equal(t, len(items), len(youtubeItems))
	for i := range items {
		our := youtubeItems[i]
		their := models.YoutubeData(items[i].LikedData)
		assert.Equal(t, their, our)

		xords = append(xords, items[i].Xord)
	}

	return items, xords
}

func TestMigration(t *testing.T) {
	x, cleanup := setupDatabase(t)
	defer cleanup()

	dumpLists(t, x)

	meta := x.liked.debugInfo().catalog.Meta
	assert.Equal(t, &models.CatalogMeta{
		ObjectIDField:    "youtube_id",
		IsUniqueObjectID: true,
	}, meta)
}

func TestSyncAppendShuffle(t *testing.T) {
	x, cleanup := setupDatabase(t)
	defer cleanup()

	liked := x.liked
	cursor := x.likesCursor

	for i := 0; i < 20; i++ {
		res, err := liked.ExecuteQuery(":sync")
		assert.NoError(t, err)

		// wait for operation to finish
		<-res.DoneChannel

		_, xords := fetchVerify(t, liked, cursor.Items)

		t.Log("xords", xords)
		t.Log("xord metric", xordMetric(xords))

		// add new item, shuffle items
		cursor.Items = append(cursor.Items, cursor.GenItem(i))
		cursor.Shuffle()
	}

	dumpLists(t, x)
}

func TestPushFrontSync(t *testing.T) {
	x, cleanup := setupDatabase(t)
	defer cleanup()

	liked := x.liked
	cursor := x.likesCursor

	for i := 0; i < 20; i++ {
		res, err := liked.ExecuteQuery(":sync")
		assert.NoError(t, err)

		// wait for operation to finish
		<-res.DoneChannel

		_, xords := fetchVerify(t, liked, cursor.Items)

		t.Log("xords", xords)
		t.Log("xord metric", xordMetric(xords))

		// append to the front
		newItem := cursor.GenItem(i)
		cursor.Items = append([]models.YoutubeData{newItem}, cursor.Items...)
	}

	dumpLists(t, x)
}

func equalItems(t *testing.T, fetched []proto.ListItem, expected []models.YoutubeData) {
	assert.Equal(t, len(fetched), len(expected))
	for i := range fetched {
		assert.Equal(t, expected[i], models.YoutubeData(fetched[i].(models.LikedModel).LikedData))
	}
}

func TestSyncManyItems(t *testing.T) {
	x, cleanup := setupDatabase(t)
	defer cleanup()

	liked := x.liked
	cursor := x.likesCursor

	const itemsCount = 1234

	for i := 0; i < itemsCount; i++ {
		cursor.Items = append(cursor.Items, cursor.GenItem(i))
	}

	res, err := liked.ExecuteQuery(":sync")
	assert.NoError(t, err)

	// wait for operation to finish
	<-res.DoneChannel

	_, xords := fetchVerify(t, liked, cursor.Items)
	t.Log("xord metric", xordMetric(xords))

	// test query functions after data has been filled
	items, err := liked.ListItems()
	assert.NoError(t, err)
	equalItems(t, items, cursor.Items)

	// test page requests
	testPage := func(limit, offset int) {
		expected := cursor.Items[offset:]
		if len(expected) > limit {
			expected = expected[:limit]
		}

		items, err := liked.PageItems(&proto.PageRequest{Limit: limit, Offset: offset})
		assert.NoError(t, err)
		equalItems(t, items, expected)
	}
	testPage(10, 0)
	testPage(itemsCount-(itemsCount%50), 50)
}
