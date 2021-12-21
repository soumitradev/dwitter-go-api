<template>
  <div class="bg-neutral-99 max-w-xl pl-4 pb-6 flex flex-row justify-between">
    <div class="flex flex-row pt-4">
      <img :src="pfpURL" class="w-12 h-12 rounded-full" />
      <div class="flex flex-col ml-4">
        <div class="flex flex-col hover:underline">
          <span class="text-left font-bold text-neutral-10">{{ name }}</span>
          <span class="text-left text-neutralVariant-50">@{{ username }}</span>
        </div>
        <span class="text-left pt-2 text-neutral-10">{{ bio }}</span>
      </div>
    </div>
    <div class="flex flex-row self-start p-2" v-if="isInUserArray(viewUser, followers)">
      <button
        class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-secondary-90 hover:text-secondary-40 transition duration-200 ease-in-out"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            d="M11 6a3 3 0 11-6 0 3 3 0 016 0zM14 17a6 6 0 00-12 0h12zM13 8a1 1 0 100 2h4a1 1 0 100-2h-4z"
          />
        </svg>
      </button>
    </div>
    <div class="flex flex-row self-start p-2" v-else>
      <button
        class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
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
</template>

<script>

export default {
  name: "UserListObject",
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
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
