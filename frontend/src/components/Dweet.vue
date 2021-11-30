<template>
  <div class="hello">
    <h1>{{ dweetBody }}</h1>
  </div>
</template>

<script>
export default {
  name: "Dweet",
  props: {
    dweetBody: String,
    id: String,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
