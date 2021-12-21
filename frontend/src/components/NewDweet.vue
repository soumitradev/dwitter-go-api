<template>
  <div class="flex flex-col bg-primary-99 max-w-xl h-fit">
    <div class="flex flex-row px-4 pt-4 pb-1">
      <img :src="viewUser.pfpURL" class="w-12 h-12 rounded-full" />
      <div class="flex flex-col w-full">
        <textarea
          class="bg-neutral-99 mx-4 mt-2 text-xl w-full border-0 resize-none appearance-none border-none outline-none"
          name="text"
          oninput="this.style.height = '';this.style.height = (this.scrollHeight + 1) + 'px'"
          rows="1"
          placeholder="What's happening?"
        ></textarea>
        <div class="flex flex-row justify-between grow mt-6 items-end">
          <div class="relative overflow-hidden inline-block self-end ml-2 group rounded-full">
            <button
              type="button"
              class="text-primary-10 bg-primary-90 border-none rounded-full transition duration-200 ease-in-out"
            >
              <div
                class="rounded-full py-2 px-4 group-hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                  />
                </svg>
              </div>
            </button>
            <input
              type="file"
              name="myfile"
              class="opacity-0 absolute left-0 top-0 text-9xl cursor-pointer"
            />
          </div>
          <div class="relative overflow-hidden inline-block self-end ml-2">
            <button
              type="button"
              class="text-primary-10 bg-primary-90 border-none rounded-full self-end mr-2 transition duration-200 ease-in-out"
            >
              <div
                class="rounded-full py-2 px-4 hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"
                  />
                </svg>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "NewDweet",
  props: {
    viewUser: {
      type: Object,
    },
  },
  components: {
  },
}
</script>



<style scoped>
.upload-btn-wrapper input[type="file"] {
  font-size: 100px;
  position: absolute;
  left: 0;
  top: 0;
  opacity: 0;
}
</style>
