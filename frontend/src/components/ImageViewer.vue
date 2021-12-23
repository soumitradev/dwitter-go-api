<template>
  <div class="flex flex-wrap justify-center">
    <div
      v-if="thumbList.length === 0"
      v-for="(item, index) in mediaList"
      :key="item.id"
      class="flex grow shrink basis-40 m-2 rounded shadow-md overflow-hidden max-h-80"
    >
      <div class="absolute m-2" v-if="editEnabled">
        <button
          type="button"
          class="text-error-10 bg-error-90 border-none rounded-full self-end transition duration-200 ease-in-out"
          @click="mediaList.splice(index, 1)"
        >
          <div
            class="rounded-full p-1 hover:bg-error-10/s2 focus:bg-error-10/s2 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 fill-current"
              viewBox="0 0 24 24"
            >
              <path d="M0 0h24v24H0V0z" fill="none" />
              <path
                d="M18.3 5.71c-.39-.39-1.02-.39-1.41 0L12 10.59 7.11 5.7c-.39-.39-1.02-.39-1.41 0-.39.39-.39 1.02 0 1.41L10.59 12 5.7 16.89c-.39.39-.39 1.02 0 1.41.39.39 1.02.39 1.41 0L12 13.41l4.89 4.89c.39.39 1.02.39 1.41 0 .39-.39.39-1.02 0-1.41L13.41 12l4.89-4.89c.38-.38.38-1.02 0-1.4z"
              />
            </svg>
          </div>
        </button>
      </div>

      <img
        v-if="item.type.startsWith('image/')"
        class="object-cover w-full block"
        :src="getURL(item)"
      />
      <video class="object-cover w-full block" v-else controls>
        <source :src="getURL(item)" :type="item.type" />
      </video>
    </div>

    <div
      v-else
      v-for="item in thumbList"
      :key="item.id"
      class="flex grow shrink basis-40 m-2 rounded shadow-md overflow-hidden max-h-80 relative justify-center items-center"
    >
      <ExternalMediaContainer :thumbURL="item.url" :type="item.type" :originalURL="item.original" />
    </div>
  </div>
</template>

<script>
import ExternalMediaContainer from "../components/ExternalMediaContainer.vue";

export default {
  name: "ImageViewer",
  methods: {
    getURL: function (file) {
      return URL.createObjectURL(file);
    },
  },
  props: {
    editEnabled: {
      type: Boolean,
    },
    mediaList: {
      type: Array,
    },
    thumbList: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  components: {
    ExternalMediaContainer,
  },
}
</script>



<style scoped>
</style>

