<script setup lang="ts">
import { computed } from 'vue'

interface Props {
    youtubeID: string,
    title: string,
    author: string,
    channelID: string,
    itemID: string,
    xord: string
}

const props = defineProps<Props>()

const channelURL = computed(() => `https://www.youtube.com/channel/${props.channelID}`)
const ariaLabel = computed(() => `${props.title} by ${props.author}`)
const watchURL = computed(() => `https://www.youtube.com/watch?v=${props.youtubeID}`)
const thumbnailURL = computed(() => `https://i.ytimg.com/vi/${props.youtubeID}/hqdefault.jpg`)
</script>

<template>
  <div class="card">
    <!--- preview can be compressed ~2 times, ?sqp=-oaymwEcCPYBEIoBSFXyq4qpAw4IARUAAIhCGAFwAcABBg==&amp;rs=AOn4CLCQUZiCTPUUhvrJvIB9bDesJgKw9w -->
    <img id="img" alt="" :src="thumbnailURL" height="100" />
    <!--- TODO: display progress, like watched 21% of video -->
    <div id="meta">
      <span>
        <a
          id="video-title"
          :href="watchURL"
          :aria-label="ariaLabel"
          :title="title"
        >
          {{ title }}
        </a>
      </span>
      <div id="metadata">
        <div id="byline-container">
          <div id="text-container">
            <a spellcheck="false" :href="channelURL" dir="auto">{{ author }}</a>
          </div>
          <div id="separator" />
        </div>
      </div>
    </div>
    <div class="internal">
      <p>youtube: {{ youtubeID }}</p>
      <p>xord: {{ xord }}</p>
      <p>id: {{ itemID }}</p>
    </div>
  </div>
</template>

<style scoped>
/* Disable link styling  */
a {
  color: inherit;
  text-decoration: none;
}

.card {
  width: 100%;
  height: 100%;
  border-radius: 5px;
  box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.16), 0 2px 10px 0 rgba(0, 0, 0, 0.12);
  display: flex;
  flex-direction: row;
  margin: 10px;
}

#meta {
  padding: 0.5rem;
}

#video-title {
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  color: #000;
}

.internal {
  margin-left: auto;
}

.internal > p {
  margin: 0.1rem;
  padding: 0;
  font-size: 0.8rem;
}
</style>
