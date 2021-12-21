<template>
  <div class="flex flex-col bg-neutral-99 max-w-xl pl-4">
    <div class="flex flex-row justify-between">
      <div class="flex flex-row">
        <img :src="pfpURL" class="w-32 h-32 rounded-full" />
        <div class="flex flex-col ml-4 text-xl self-center">
          <span class="text-left font-bold text-neutral-10">{{ name }}</span>
          <span class="text-left text-neutralVariant-50">@{{ username }}</span>
        </div>
      </div>
      <div class="flex flex-row self-start p-4" v-if="isInUserArray(viewUser, followers)">
        <button
          class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-secondary-90 hover:text-secondary-40 transition duration-200 ease-in-out"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              d="M11 6a3 3 0 11-6 0 3 3 0 016 0zM14 17a6 6 0 00-12 0h12zM13 8a1 1 0 100 2h4a1 1 0 100-2h-4z"
            />
          </svg>
        </button>
      </div>
      <div class="flex flex-row self-start p-4" v-else>
        <button
          class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              d="M8 9a3 3 0 100-6 3 3 0 000 6zM8 11a6 6 0 016 6H2a6 6 0 016-6zM16 7a1 1 0 10-2 0v1h-1a1 1 0 100 2h1v1a1 1 0 102 0v-1h1a1 1 0 100-2h-1V7z"
            />
          </svg>
        </button>
      </div>
    </div>

    <span class="text-left px-2 pt-6 text-neutral-10">{{ bio }}</span>

    <span :if="email" class="text-left px-2 pt-1 text-neutralVariant-50">{{ email }}</span>
    <div class="flex flex-row text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5 mx-1"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fill-rule="evenodd"
          d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
          clip-rule="evenodd"
        />
      </svg>
      <span>Joined {{ formatDate(createdAt) }}</span>
    </div>

    <div class="flex flex-row text-left py-3 text-neutralVariant-50">
      <div class="hover:underline">
        <span class="font-bold px-1 text-neutral-40">{{ followerCount }}</span>
        <span>{{ followerCount == 1 ? "Follower" : "Followers" }}</span>
      </div>
      <div class="hover:underline">
        <span class="ml-4 font-bold px-1 text-neutral-40">{{ followingCount }}</span>
        <span>{{ followingCount == 1 ? "Following" : "Following" }}</span>
      </div>
    </div>
    <div class="flex flex-row mb-2 text-left text-sm text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-4 w-4 mr-2"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"
        />
      </svg>
      <span v-if="following.length > 0">Follows {{ overflowNames(following) }}</span>
      <span v-else>Doesn't follow anyone you follow</span>
    </div>
    <div class="flex flex-row mb-2 text-left text-sm text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-4 w-4 mr-2"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"
        />
      </svg>
      <span v-if="followers.length > 0">Followed by {{ overflowNames(followers) }}</span>
      <span v-else>Isn't followed by anyone you follow</span>
    </div>
  </div>
</template>

<script>
import { Menu, MenuButton, MenuItems, MenuItem } from "@headlessui/vue";

export default {
  name: "User",
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
    isInUserArray: function (item, array) {
      return array.map(e => e.username).includes(item);
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
    username: {
      type: String,
    },
    name: {
      type: String,
    },
    email: {
      type: String,
    },
    bio: {
      type: String,
    },
    pfpURL: {
      type: String,
    },
    followers: {
      type: Array,
    },
    followerCount: {
      type: Number,
    },
    following: {
      type: Array,
    },
    followingCount: {
      type: Number,
    },
    createdAt: {
      type: String,
    },
    viewUser: {
      type: String,
    }
  },
  components: {
    Menu,
    MenuButton,
    MenuItems,
    MenuItem,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
