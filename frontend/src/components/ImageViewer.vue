<template>
  <div class="flex flex-wrap justify-center">
    <div
      v-for="(item, index) in mediaList"
      :key="item.id"
      class="flex grow shrink basis-40 m-2 rounded shadow-md overflow-hidden max-h-80"
    >
      <div class="absolute m-2" v-if="editEnabled">
        <button
          type="button"
          class="text-error-10 bg-error-90 border-none rounded-full self-end mr-2 transition duration-200 ease-in-out"
          @click="mediaList.splice(index, 1)"
        >
          <div
            class="rounded-full p-1 hover:bg-error-10/s2 focus:bg-error-10/s2 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clip-rule="evenodd"
              />
            </svg>
          </div>
        </button>
      </div>

      <img
        v-if="item.type.startsWith('image/')"
        class="object-cover w-full block"
        :src="getURL(item)"
        alt
      />
      <video class="object-cover w-full block" v-else controls>
        <source :src="getURL(item)" :type="item.type" />
      </video>
    </div>
  </div>
</template>

<script>

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
  },
  components: {
  },
}
</script>



<style scoped>
</style>
