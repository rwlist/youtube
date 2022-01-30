<script setup lang="ts">
import { toRefs } from 'vue'
import { useListCtl } from '../composables/useListCtl'
import YoutubeCard from './YoutubeCard.vue'

const props = defineProps<{ listId: string }>()

const { listId } = toRefs(props)

const { status, allItems, query, executeQuery } = useListCtl(listId)
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
    <div class="list-items">
      <YoutubeCard v-for="item in allItems" :key="item.ItemID" v-bind="item" />
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
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  padding: 10px;
}
</style>
