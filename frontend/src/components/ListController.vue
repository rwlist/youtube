<script setup lang="ts">
import { toRefs } from 'vue'
import { useListCtl } from '../composables/useListCtl'
import type { ItemLiked } from '../rpc/proto_gen'
import YoutubeCard from './YoutubeCard.vue'
import Grid from 'vue-virtual-scroll-grid'

const props = defineProps<{ listId: string }>()

// TODO: don't always cast to ItemLiked

const placeholderProp: ItemLiked = {
    YoutubeID: '',
    Title: '',
    Author: '',
    ChannelID: '',
    ItemID: 0,
    Xord: '',
}

const { listId } = toRefs(props)

const { status, allItems, query, executeQuery, supportsPages, getPagedList } =
    useListCtl(listId)

const pageProvider = (
    pageNumber: number,
    pageSize: number,
): Promise<ItemLiked[]> => {
    return getPagedList()
        .fetchPage(pageNumber * pageSize, pageSize)
        .then((listItems) => {
            const it = listItems as ItemLiked[]

            if (it.length < pageSize) {
                return [
                    ...it,
                    ...(Array(pageSize - it.length).fill(
                        placeholderProp,
                    ) as ItemLiked[]),
                ]
            }
            return it
        })
}
</script>

<template>
    <div class="list-controller">
        <div class="ctrl-panel">
            <!-- Search query -->
            <div class="search-query">
                <input
                    v-model="query"
                    class="search-query-input"
                    type="text"
                    placeholder="Enter query"
                    @keyup.enter="executeQuery"
                />
            </div>
            <br />
            <!-- Control/status pane -->
            <fieldset>
                <legend>{{ status.header }}</legend>
                <div class="ctrl-panel-status">
                    <!-- display response as a monospace text -->

                    <code>
                        <pre>{{ status.response }}</pre>
                    </code>

                    <!-- <div class="ctrl-panel-status-item">
            <span class="ctrl-panel-status-item-label"> Total: </span>
            <span class="ctrl-panel-status-item-value"> 123 </span>
          </div>
          <div class="ctrl-panel-status-item">
            <span class="ctrl-panel-status-item-label"> Selected: </span>
            <span class="ctrl-panel-status-item-value"> 21412 </span>
          </div> -->
                </div>
            </fieldset>
        </div>

        <div v-if="supportsPages">
            <Grid
                :length="getPagedList().getInfo().ItemsCount"
                :page-size="40"
                :page-provider="pageProvider"
                class="list-items"
            >
                <template #probe="{ style }">
                    <YoutubeCard :style="style" v-bind="placeholderProp" />
                </template>

                <!-- When the item is not loaded, a placeholder is rendered -->
                <template #placeholder="{ style }">
                    <YoutubeCard :style="style" v-bind="placeholderProp" />
                </template>

                <!-- Render a loaded item -->
                <template #default="{ item, style }">
                    <YoutubeCard :style="style" v-bind="item" />
                </template>
            </Grid>
        </div>
        <div v-else class="list-items">
            <YoutubeCard
                v-for="item in allItems"
                :key="(item as ItemLiked).ItemID"
                v-bind="(item as ItemLiked)"
            />
        </div>
    </div>
</template>

<style scoped>
.search-query-input {
    font-size: 1.3rem;
    border: 1px solid #e7e7e7;
    border-radius: 3px;
    padding: 10px;
    width: 100%;
}

/* center content horizontally and limit width */
.list-controller {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100%;
}

.list-controller > * {
    width: 70%;
}

.list-items {
    display: grid;
    padding: 10px;
    grid-gap: 10px;
}
</style>
