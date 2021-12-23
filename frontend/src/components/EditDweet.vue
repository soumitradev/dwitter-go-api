<template>
  <ReplyParentDweet v-if="isReply" v-bind="replyTo" :viewUser="viewUser" />
  <div
    class="flex flex-col bg-neutral-99 max-w-xl divide-y divide-neutralVariant-60 divide-opacity-20"
    :class="{ 'px-4 pb-4': isReply, 'p-4': !isReply }"
  >
    <div>
      <div
        v-if="likeUsers.length > 0"
        class="text-left text-sm flex flex-col mb-1 text-neutralVariant-50"
      >
        <div class="flex flex-row mb-1">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            class="h-4 w-4 mr-2 fill-current"
          >
            <path d="M0 0h24v24H0V0z" fill="none" />
            <path
              d="M13.35 20.13c-.76.69-1.93.69-2.69-.01l-.11-.1C5.3 15.27 1.87 12.16 2 8.28c.06-1.7.93-3.33 2.34-4.29 2.64-1.8 5.9-.96 7.66 1.1 1.76-2.06 5.02-2.91 7.66-1.1 1.41.96 2.28 2.59 2.34 4.29.14 3.88-3.3 6.99-8.55 11.76l-.1.09z"
            />
          </svg>
          <span>{{ overflowNames(likeUsers) }} liked</span>
        </div>
      </div>
      <div
        v-if="redweetUsers.length > 0"
        class="text-left text-sm flex flex-col mb-2 text-neutralVariant-50"
      >
        <div class="flex flex-row mb-1">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            class="h-4 w-4 mr-2 fill-current"
          >
            <path
              d="M7 7h10v1.79c0 .45.54.67.85.35l2.79-2.79c.2-.2.2-.51 0-.71l-2.79-2.79c-.31-.31-.85-.09-.85.36V5H6c-.55 0-1 .45-1 1v4c0 .55.45 1 1 1s1-.45 1-1V7zm10 10H7v-1.79c0-.45-.54-.67-.85-.35l-2.79 2.79c-.2.2-.2.51 0 .71l2.79 2.79c.31.31.85.09.85-.36V19h11c.55 0 1-.45 1-1v-4c0-.55-.45-1-1-1s-1 .45-1 1v3z"
            />
          </svg>
          <span>{{ overflowNames(redweetUsers) }} redweeted</span>
        </div>
      </div>
      <div class="flex justify-between">
        <div class="flex flex-row">
          <img :src="author.pfpURL" class="w-12 h-12 rounded-full" />
          <div class="flex flex-col ml-4 hover:underline">
            <span class="text-left font-bold text-neutral-10">{{ author.name }}</span>
            <span class="text-left text-neutralVariant-50">@{{ author.username }}</span>
          </div>
        </div>
        <div>
          <Menu as="div" class="relative inline-block text-left">
            <div>
              <MenuButton
                class="h-10 w-10 inline-flex justify-center content-center text-sm rounded-full font-medium text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out focus-visible:ring-2 focus-visible:ring-neutralVariant-50 focus-visible:ring-opacity-75"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6 place-self-center fill-current"
                  viewBox="0 0 24 24"
                >
                  <path d="M0 0h24v24H0V0z" fill="none" />
                  <path
                    d="M6 10c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm12 0c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm-6 0c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"
                  />
                </svg>
              </MenuButton>
            </div>

            <transition
              enter-active-class="transition duration-100 ease-out"
              enter-from-class="transform scale-95 opacity-0"
              enter-to-class="transform scale-100 opacity-100"
              leave-active-class="transition duration-75 ease-in"
              leave-from-class="transform scale-100 opacity-100"
              leave-to-class="transform scale-95 opacity-0"
            >
              <MenuItems
                class="z-10 absolute right-0 w-56 mt-2 origin-top-right bg-neutral-99 divide-y divide-neutralVariant-50 rounded-md shadow-lg ring-1 ring-neutralVariant-50 ring-opacity-5"
              >
                <div class="px-1 py-1">
                  <MenuItem v-slot="{ active }" class="transition duration-200 ease-in-out">
                    <button
                      :class="[
                        active
                          ? 'bg-neutralVariant-90 text-neutralVariant-30'
                          : 'text-neutral-10',
                        'group flex rounded-md items-center w-full px-2 py-2 text-sm',
                      ]"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 mr-2 fill-current"
                        viewBox="0 0 24 24"
                      >
                        <g>
                          <rect fill="none" height="24" width="24" />
                          <rect fill="none" height="24" width="24" />
                        </g>
                        <g>
                          <path
                            d="M14,8c0-2.21-1.79-4-4-4S6,5.79,6,8s1.79,4,4,4S14,10.21,14,8z M2,18v1c0,0.55,0.45,1,1,1h14c0.55,0,1-0.45,1-1v-1 c0-2.66-5.33-4-8-4S2,15.34,2,18z M18,10h4c0.55,0,1,0.45,1,1v0c0,0.55-0.45,1-1,1h-4c-0.55,0-1-0.45-1-1v0 C17,10.45,17.45,10,18,10z"
                          />
                        </g>
                      </svg>
                      Unfollow {{ author.name }}
                    </button>
                  </MenuItem>
                  <MenuItem
                    v-if="author.username == viewUser"
                    v-slot="{ active }"
                    class="transition duration-200 ease-in-out"
                  >
                    <button
                      :class="[
                        active
                          ? 'bg-neutralVariant-90 text-neutralVariant-30'
                          : 'text-neutral-10',
                        'group flex rounded-md items-center w-full px-2 py-2 text-sm',
                      ]"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 mr-2 fill-current"
                        viewBox="0 0 24 24"
                      >
                        <path d="M0 0h24v24H0V0z" fill="none" />
                        <path
                          d="M3 17.46v3.04c0 .28.22.5.5.5h3.04c.13 0 .26-.05.35-.15L17.81 9.94l-3.75-3.75L3.15 17.1c-.1.1-.15.22-.15.36zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z"
                        />
                      </svg>
                      Edit Dweet
                    </button>
                  </MenuItem>
                  <MenuItem
                    v-if="author.username == viewUser"
                    v-slot="{ active }"
                    class="transition duration-200 ease-in-out"
                  >
                    <button
                      :class="[
                        active
                          ? 'bg-error-90 text-error-10'
                          : 'text-neutral-10',
                        'group flex rounded-md items-center w-full px-2 py-2 text-sm',
                      ]"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5 mr-2 fill-current"
                        viewBox="0 0 24 24"
                      >
                        <path d="M0 0h24v24H0V0z" fill="none" />
                        <path
                          d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V9c0-1.1-.9-2-2-2H8c-1.1 0-2 .9-2 2v10zM18 4h-2.5l-.71-.71c-.18-.18-.44-.29-.7-.29H9.91c-.26 0-.52.11-.7.29L8.5 4H6c-.55 0-1 .45-1 1s.45 1 1 1h12c.55 0 1-.45 1-1s-.45-1-1-1z"
                        />
                      </svg>
                      Delete Dweet
                    </button>
                  </MenuItem>
                </div>
              </MenuItems>
            </transition>
          </Menu>
        </div>
        <!-- <button class="p-4">...</button> -->
      </div>
      <div class="text-left my-4">
        <textarea
          class="bg-neutral-99 mx-4 mt-2 text-xl w-full border-0 resize-none appearance-none border-none outline-none"
          name="text"
          oninput="this.style.height = '';this.style.height = (this.scrollHeight + 1) + 'px'"
          rows="1"
          :value="dweetBody"
          autofocus
        ></textarea>
      </div>

      <ImageViewer class="m-4" :editEnabled="true" :mediaList="fileList" />

      <div class="flex flex-row justify-between grow mt-2 mx-4 items-end">
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

      <div class="flex flex-col text-left my-4 mt-auto text-neutralVariant-50">
        <span class="text-sm">
          <b>Posted:</b>
          {{ formatDate(postedAt) }}
          <b class="mx-1">·</b>
          <b>Last Updated:</b>
          {{ formatDate(lastUpdatedAt) }}
        </span>
        <span class="text-sm"></span>
      </div>
    </div>
    <div class="flex flex-row text-left py-3 text-neutralVariant-50">
      <div class="hover:underline">
        <span class="font-bold px-1 text-neutral-40">{{ redweetCount }}</span>
        <span>{{ redweetCount == 1 ? "Redweet" : "Redweets" }}</span>
      </div>
      <div class="hover:underline">
        <span class="ml-4 font-bold px-1 text-neutral-40">{{ likeCount }}</span>
        <span>{{ likeCount == 1 ? "Like" : "Likes" }}</span>
      </div>
    </div>
    <div>
      <div class="flex flex-row justify-between pt-1">
        <div class="grow">
          <button
            class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-current"
              viewBox="0 0 24 24"
            >
              <path d="M0 0h24v24H0V0z" fill="none" />
              <path
                d="M10 9V7.41c0-.89-1.08-1.34-1.71-.71L3.7 11.29c-.39.39-.39 1.02 0 1.41l4.59 4.59c.63.63 1.71.19 1.71-.7V14.9c5 0 8.5 1.6 11 5.1-1-5-4-10-11-11z"
              />
            </svg>
          </button>
        </div>
        <div class="grow">
          <button
            class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-tertiary-90 hover:text-tertiary-40 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-tertiary-60"
              viewBox="0 0 24 24"
              v-if="redweetUsers.map(x => x.username).includes(viewUser)"
            >
              <path
                d="M7 7h10v1.79c0 .45.54.67.85.35l2.79-2.79c.2-.2.2-.51 0-.71l-2.79-2.79c-.31-.31-.85-.09-.85.36V5H6c-.55 0-1 .45-1 1v4c0 .55.45 1 1 1s1-.45 1-1V7zm10 10H7v-1.79c0-.45-.54-.67-.85-.35l-2.79 2.79c-.2.2-.2.51 0 .71l2.79 2.79c.31.31.85.09.85-.36V19h11c.55 0 1-.45 1-1v-4c0-.55-.45-1-1-1s-1 .45-1 1v3z"
              />
            </svg>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-current"
              viewBox="0 0 24 24"
              v-else
            >
              <path
                d="M7 7h10v1.79c0 .45.54.67.85.35l2.79-2.79c.2-.2.2-.51 0-.71l-2.79-2.79c-.31-.31-.85-.09-.85.36V5H6c-.55 0-1 .45-1 1v4c0 .55.45 1 1 1s1-.45 1-1V7zm10 10H7v-1.79c0-.45-.54-.67-.85-.35l-2.79 2.79c-.2.2-.2.51 0 .71l2.79 2.79c.31.31.85.09.85-.36V19h11c.55 0 1-.45 1-1v-4c0-.55-.45-1-1-1s-1 .45-1 1v3z"
              />
            </svg>
          </button>
        </div>
        <div class="grow">
          <button
            class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-secondary-90 hover:text-secondary-40 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-secondary-50"
              viewBox="0 0 24 24"
              v-if="likeUsers.map(x => x.username).includes(viewUser)"
            >
              <path
                d="M13.35 20.13c-.76.69-1.93.69-2.69-.01l-.11-.1C5.3 15.27 1.87 12.16 2 8.28c.06-1.7.93-3.33 2.34-4.29 2.64-1.8 5.9-.96 7.66 1.1 1.76-2.06 5.02-2.91 7.66-1.1 1.41.96 2.28 2.59 2.34 4.29.14 3.88-3.3 6.99-8.55 11.76l-.1.09z"
              />
            </svg>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-current"
              viewBox="0 0 24 24"
              v-else
            >
              <path
                d="M19.66 3.99c-2.64-1.8-5.9-.96-7.66 1.1-1.76-2.06-5.02-2.91-7.66-1.1-1.4.96-2.28 2.58-2.34 4.29-.14 3.88 3.3 6.99 8.55 11.76l.1.09c.76.69 1.93.69 2.69-.01l.11-.1c5.25-4.76 8.68-7.87 8.55-11.75-.06-1.7-.94-3.32-2.34-4.28zM12.1 18.55l-.1.1-.1-.1C7.14 14.24 4 11.39 4 8.5 4 6.5 5.5 5 7.5 5c1.54 0 3.04.99 3.57 2.36h1.87C13.46 5.99 14.96 5 16.5 5c2 0 3.5 1.5 3.5 3.5 0 2.89-3.14 5.74-7.9 10.05z"
              />
            </svg>
          </button>
        </div>
        <div class="grow">
          <button
            class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="2-6 h-6 fill-current"
              viewBox="0 0 24 24"
            >
              <path d="M0 0h24v24H0V0z" fill="none" />
              <path
                d="M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.5-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92s2.92-1.31 2.92-2.92-1.31-2.92-2.92-2.92z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import ImageViewer from "../components/ImageViewer.vue";
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
import ReplyParentDweet from "../components/ReplyParentDweet.vue";

export default {
  name: "Dweet",
  methods: {
    formatDate: function (date) {
      let dateObj = new Date(date);
      return (
        dateObj.toLocaleString("en-US", {
          hour12: true,
          timeStyle: "short",
        }) +
        " · " +
        dateObj.toLocaleString("en-US", {
          dateStyle: "medium",
        })
      );
    },
    overflowNames: function (list) {
      // TODO: Remove viewer's name from list
      if (list.length > 0) {
        let text = "";
        let i = 0;
        let curLen = 0;
        while (curLen <= 42 && i < list.length) {
          curLen += list[i].name.length + 2;
          i++;
        }

        let j = 0;
        let numberOfNames = list.length - i;

        while (i > 0) {
          if (i == 1) {
            if (list.length > 1) {
              text += "and " + list[j].name;
            } else {
              text += list[j].name;
            }
          } else {
            text += list[j].name + ", ";
          }
          i--;
          j++;
        }

        if (numberOfNames == 0) {
          return text;
        } else if (numberOfNames > 1) {
          text += " and " + numberOfNames + " others";
        } else if (numberOfNames == 1) {
          text += " and " + numberOfNames + " other";
        }
        return text;
      } else {
        return "";
      }
    },
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
    dweetBody: {
      type: String,
    },
    id: {
      type: String,
    },
    author: {
      type: Object,
    },
    authorID: {
      type: String,
    },
    postedAt: {
      type: String,
    },
    lastUpdatedAt: {
      type: String,
    },
    likeCount: {
      type: Number,
    },
    likeUsers: {
      type: Array,
    },
    isReply: {
      type: Boolean,
    },
    originalReplyID: {
      type: String,
    },
    replyTo: {
      type: Object,
    },
    replyCount: {
      type: Number,
    },
    replyDweets: {
      type: Array,
    },
    redweetCount: {
      type: Number,
    },
    redweetUsers: {
      type: Array,
    },
    media: {
      type: Array,
    },
    viewUser: {
      type: String,
    },
  },
  components: {
    Menu,
    MenuButton,
    MenuItems,
    MenuItem,
    ReplyParentDweet,
    ImageViewer,
  },
  setup() {
    const fileList = ref([]);


    return {
      fileList,
    }
  },
  created() {
    for (let i = 0; i < this.media.length; i++) {
      let type = "";
      switch (this.media[i].replace(/\?.+/gi, "").split(".").pop().toLowerCase()) {
        case "jpg":
          type = "image/jpeg"
          break;
        case "jpeg":
          type = "image/jpeg"
          break;
        case "gif":
          type = "image/gif"
          break;
        case "png":
          type = "image/png"
          break;
        case "mp4":
          type = "video/mp4"
          break;
      }

      fetch(this.media[i]).then(
        r => r.blob()
      ).then(
        blob => new File([blob], this.media[i].replace(/\?.+/gi, "").split('media%2F')[1], { type: type })
      ).then(
        file => this.fileList.push(file)
      );
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
