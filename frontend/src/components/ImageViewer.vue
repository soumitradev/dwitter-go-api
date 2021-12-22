<template>
  <div class="flex flex-wrap justify-center">
    <div
      v-for="item in mediaList"
      :key="item.id"
      class="flex grow shrink basis-40 m-2 rounded shadow-md overflow-hidden max-h-80"
    >
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
