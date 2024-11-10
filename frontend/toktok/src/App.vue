<template>
  <div id="app">
    <h1>Video Player TikTok</h1>
    <VideoPlayer :autoplay="true" />
  </div>
</template>

<script>
import VideoPlayer from './components/VideoPlayer.vue';

export default {
  name: 'App',
  components: {
    VideoPlayer,
  },
};
</script>

<style>
#app {
  text-align: center;
  padding: 2rem;
}
</style>
