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
            class="h-4 w-4 mr-2"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
              clip-rule="evenodd"
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
            class="h-4 w-4 mr-2"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M18 3a1 1 0 00-1.447-.894L8.763 6H5a3 3 0 000 6h.28l1.771 5.316A1 1 0 008 18h1a1 1 0 001-1v-4.382l6.553 3.276A1 1 0 0018 15V3z"
              clip-rule="evenodd"
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
                  class="h-5 w-5 place-self-center"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z"
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
                class="absolute right-0 w-56 mt-2 origin-top-right bg-neutral-99 divide-y divide-neutralVariant-50 rounded-md shadow-lg ring-1 ring-neutralVariant-50 ring-opacity-5"
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
                        class="h-5 w-5 mr-2"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          d="M11 6a3 3 0 11-6 0 3 3 0 016 0zM14 17a6 6 0 00-12 0h12zM13 8a1 1 0 100 2h4a1 1 0 100-2h-4z"
                        />
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
                        class="h-5 w-5 mr-2"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
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
                        class="h-5 w-5 mr-2"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                          clip-rule="evenodd"
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
        <span class="text-2xl lea text-neutral-10">{{ dweetBody }}</span>
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
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
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
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"
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
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
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
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";
import ReplyParentDweet from "../components/ReplyParentDweet.vue";

export default {
  name: "Dweet",
  methods: {
    formatDate: function (date) {
      var dateObj = new Date(date);
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
        var text = "";
        var i = 0;
        var curLen = 0;
        while (curLen <= 42 && i < list.length) {
          curLen += list[i].name.length + 2;
          i++;
        }

        var j = 0;
        var numberOfNames = list.length - i;

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
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
