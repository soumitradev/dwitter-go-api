<template>
  <div class="absolute" v-if="(type === 'video/mp4') && (!this.active)">
    <button
      type="button"
      class="text-primary-10 bg-primary-90 border-none rounded-full self-end transition duration-200 ease-in-out opacity-60"
      @click="activateVideo()"
    >
      <div
        class="rounded-full p-1 hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          enable-background="new 0 0 24 24"
          class="h-8 w-8"
          viewBox="0 0 24 24"
          fill="#000000"
        >
          <path
            d="M8 6.82v10.36c0 .79.87 1.27 1.54.84l8.14-5.18c.62-.39.62-1.29 0-1.69L9.54 5.98C8.87 5.55 8 6.03 8 6.82z"
          />
        </svg>
      </div>
    </button>
  </div>
  <div class="absolute" v-if="(type === 'image/gif')">
    <button
      type="button"
      class="text-primary-10 bg-primary-90 border-none rounded-full self-end transition duration-200 ease-in-out opacity-60"
    >
      <a :href="originalURL" target="_blank">
        <div
          class="rounded-full p-1 hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            enable-background="new 0 0 24 24"
            class="h-8 w-8"
            viewBox="0 0 24 24"
            fill="#000000"
          >
            <path d="M0 0h24v24H0V0z" fill="none" />
            <path
              d="M12.25 9c.41 0 .75.34.75.75v4.5c0 .41-.34.75-.75.75s-.75-.34-.75-.75v-4.5c0-.41.34-.75.75-.75zM10 9.75c0-.41-.34-.75-.75-.75H6c-.6 0-1 .5-1 1v4c0 .5.4 1 1 1h3c.6 0 1-.5 1-1v-1.25c0-.41-.34-.75-.75-.75s-.75.34-.75.75v.75h-2v-3h2.75c.41 0 .75-.34.75-.75zm9 0c0-.41-.34-.75-.75-.75H15.5c-.55 0-1 .45-1 1v4.25c0 .41.34.75.75.75s.75-.34.75-.75V13h1.25c.41 0 .75-.34.75-.75s-.34-.75-.75-.75H16v-1h2.25c.41 0 .75-.34.75-.75z"
            />
          </svg>
        </div>
      </a>
    </button>
  </div>
  <div
    class="absolute"
    v-if="(type === 'image/png') || (type === 'image/jpg') || (type === 'image/jpeg')"
  >
    <button
      type="button"
      class="text-primary-10 bg-primary-90 border-none rounded-full self-end transition duration-200 ease-in-out opacity-60"
    >
      <a :href="originalURL" target="_blank">
        <div
          class="rounded-full p-2 hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            enable-background="new 0 0 24 24"
            class="h-6 w-6"
            viewBox="0 0 24 24"
            fill="#000000"
          >
            <path d="M0 0h24v24H0V0z" fill="none" />
            <path
              d="M18 19H6c-.55 0-1-.45-1-1V6c0-.55.45-1 1-1h5c.55 0 1-.45 1-1s-.45-1-1-1H5c-1.11 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-6c0-.55-.45-1-1-1s-1 .45-1 1v5c0 .55-.45 1-1 1zM14 4c0 .55.45 1 1 1h2.59l-9.13 9.13c-.39.39-.39 1.02 0 1.41.39.39 1.02.39 1.41 0L19 6.41V9c0 .55.45 1 1 1s1-.45 1-1V4c0-.55-.45-1-1-1h-5c-.55 0-1 .45-1 1z"
            />
          </svg>
        </div>
      </a>
    </button>
  </div>

  <img class="object-cover w-full block" :src="thumbURL" v-if="!this.active" />
  <video class="object-cover w-full block" v-if="this.active" controls preload="none" autoplay>
    <source :src="originalURL" :type="type" />
  </video>
</template>

<script>

export default {
  name: "ExternalMediaContainer",
  methods: {
    activateVideo: function () {
      this.active = true;
    }
  },
  props: {
    thumbURL: {
      type: String,
    },
    type: {
      type: String,
    },
    originalURL: {
      type: String,
    },
  },
  components: {
  },
  data() {
    return {
      active: false,
    }
  }
}
</script>




<style scoped>
</style>

