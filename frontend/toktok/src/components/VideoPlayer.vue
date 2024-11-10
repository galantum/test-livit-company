<template>
    <div class="video-player-container">
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Cari video..."
        @keydown.enter="fetchVideos"
      />
  
      <div v-if="loading" class="loader">Memuat...</div>
  
      <div v-if="!loading && videos.length === 0" class="no-results">
        Tidak ada video ditemukan.
      </div>
  
      <div v-for="(video, index) in videos" :key="index" class="video-wrapper">
        <video
          ref="videoRefs"
          :src="video.link"
          controls
          :muted="isMuted"
          @play="handlePlay(index)"
          @pause="handlePause"
        ></video>
  
        <div class="controls">
          <button @click="togglePlayPause(index)">
            {{ isPlaying[index] ? 'Jeda' : 'Putar' }}
          </button>
          <button @click="toggleMute">Volume {{ isMuted ? 'Mati' : 'Hidup' }}</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, watch } from 'vue';
  import { searchVideos } from '@/services/VideoService';
  
  export default {
    name: 'VideoPlayer',
    props: {
      autoplay: { type: Boolean, default: false },
    },
    setup() {
      const searchQuery = ref('');
      const videos = ref([]);
      const loading = ref(false);
      const isPlaying = ref([]);
      const isMuted = ref(true);
      const videoRefs = ref([]);
  
      const fetchVideos = async () => {
        loading.value = true;
        videos.value = await searchVideos(searchQuery.value);
        loading.value = false;
        isPlaying.value = videos.value.map(() => false); // Reset isPlaying array
      };
  
      const togglePlayPause = (index) => {
        const video = videoRefs.value[index];
        if (!video) return;
        if (video.paused) {
          video.play();
          isPlaying.value[index] = true;
        } else {
          video.pause();
          isPlaying.value[index] = false;
        }
      };
  
      const handlePlay = (index) => {
        isPlaying.value[index] = true;
      };
  
      const handlePause = () => {
        isPlaying.value.fill(false);
      };
  
      const toggleMute = () => {
        isMuted.value = !isMuted.value;
        videoRefs.value.forEach((video) => {
          if (video) video.muted = isMuted.value;
        });
      };
  
      watch(searchQuery, fetchVideos);
  
      onMounted(() => {
        if (videos.value.length === 0) fetchVideos();
      });
  
      return {
        searchQuery,
        videos,
        loading,
        videoRefs,
        isPlaying,
        isMuted,
        togglePlayPause,
        toggleMute,
        handlePlay,
        handlePause,
        fetchVideos,
      };
    },
  };
  </script>
  
  <style scoped>
  .video-player-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  .video-wrapper {
    width: 100%;
    max-width: 600px;
    margin-bottom: 2rem;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .controls {
    display: flex;
    gap: 1rem;
    margin-top: 0.5rem;
  }
  .loader, .no-results {
    text-align: center;
    font-size: 1.2rem;
  }
  @media (max-width: 768px) {
    .video-wrapper {
      max-width: 100%;
    }
  }
  </style>
  