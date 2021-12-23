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
        <ImageViewer class="m-4" :editEnabled="true" :mediaList="fileList" />
        <div class="flex flex-row justify-between grow mt-2 items-end">
          <div class="relative overflow-hidden inline-block self-end ml-2 group rounded-full">
            <button
              type="button"
              class="text-primary-10 bg-primary-90 border-none rounded-full transition duration-200 ease-in-out"
            >
              <div
                class="rounded-full p-2 group-hover:bg-primary-10/s2 focus:bg-primary-10/s2 transition duration-200 ease-in-out"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 fill-current"
                  viewBox="0 0 24 24"
                >
                  <path d="M0 0h24v24H0V0z" fill="none" />
                  <path
                    d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.9 13.98l2.1 2.53 3.1-3.99c.2-.26.6-.26.8.01l3.51 4.68c.25.33.01.8-.4.8H6.02c-.42 0-.65-.48-.39-.81L8.12 14c.19-.26.57-.27.78-.02z"
                  />
                </svg>
              </div>
            </button>
            <input
              type="file"
              class="opacity-0 absolute left-0 top-0 text-9xl cursor-pointer"
              @change="updateFiles"
              accept="image/png, image/jpeg, image/gif, video/mp4"
              multiple
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
                  class="h-6 w-6 fill-current"
                  viewBox="0 0 24 24"
                >
                  <path d="M0 0h24v24H0V0z" fill="none" />
                  <path
                    d="M3.4 20.4l17.45-7.48c.81-.35.81-1.49 0-1.84L3.4 3.6c-.66-.29-1.39.2-1.39.91L2 9.12c0 .5.37.93.87.99L17 12 2.87 13.88c-.5.07-.87.5-.87 1l.01 4.61c0 .71.73 1.2 1.39.91z"
                  />
                </svg>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="closeModal">
      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="min-h-screen px-4 text-center">
          <TransitionChild
            as="template"
            enter="duration-200 ease-out"
            enter-from="opacity-0"
            enter-to="opacity-100"
            leave="duration-100 ease-in"
            leave-from="opacity-100"
            leave-to="opacity-0"
          >
            <DialogOverlay class="fixed inset-0 bg-neutral-30/s5" />
          </TransitionChild>

          <span class="inline-block h-screen align-middle" aria-hidden="true">&#8203;</span>

          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <div
              class="inline-block w-full max-w-md p-6 my-8 overflow-hidden text-left align-middle transition-all transform bg-neutral-100 shadow-xl rounded-2xl"
            >
              <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">Error</DialogTitle>
              <div class="mt-2">
                <p class="text-sm text-gray-500">File too large (Limit is 8MB per file)</p>
              </div>

              <div class="mt-4">
                <button
                  type="button"
                  class="inline-flex justify-center text-sm font-medium text-error-10 bg-error-90 border border-transparent rounded-md focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-error-40"
                  @click="closeModal"
                >
                  <div
                    class="rounded-full py-2 px-4 group-hover:bg-error-10/s2 focus:bg-error-10/s2 transition duration-200 ease-in-out"
                  >OK</div>
                </button>
              </div>
            </div>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script>
import { ref } from 'vue'
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogOverlay,
  DialogTitle,
} from '@headlessui/vue'
import ImageViewer from "../components/ImageViewer.vue";

export default {
  name: "NewDweet",
  methods: {
    updateFiles: function (event) {
      let valid = this.validateFiles(event.target.files);
      if (valid) {
        this.fileList.push(...event.target.files);
      } else {
        this.openModal();
      }
    },
    validateFiles: function (fileList) {
      for (let fileIndex = 0; fileIndex < fileList.length; fileIndex++) {
        let sizeValid = (fileList[fileIndex].size <= (8 << 20));

        if (!sizeValid) {
          return false;
        }
      }
      return true;
    },
  },
  props: {
    viewUser: {
      type: Object,
    },
  },
  components: {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogOverlay,
    DialogTitle,
    ImageViewer,
  },
  setup() {
    const isOpen = ref(false);
    const fileList = ref([]);

    return {
      isOpen,
      fileList,
      closeModal() {
        isOpen.value = false;
      },
      openModal() {
        isOpen.value = true;
      },
    }
  },
}
</script>



<style scoped>
</style>
