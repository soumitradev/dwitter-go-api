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
    <div
      class="flex flex-row self-start p-2"
      v-if="followers.map(e => e.username).includes(viewUser)"
    >
      <button
        class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-secondary-90 hover:text-secondary-40 transition duration-200 ease-in-out"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 fill-current" viewBox="0 0 24 24">
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
      </button>
    </div>
    <div class="flex flex-row self-start p-2" v-else>
      <button
        class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 fill-current" viewBox="0 0 24 24">
          <g>
            <rect fill="none" height="24" width="24" />
          </g>
          <g>
            <path
              d="M15.39,14.56C13.71,13.7,11.53,13,9,13c-2.53,0-4.71,0.7-6.39,1.56C1.61,15.07,1,16.1,1,17.22V20h16v-2.78 C17,16.1,16.39,15.07,15.39,14.56z M9,12c2.21,0,4-1.79,4-4c0-2.21-1.79-4-4-4S5,5.79,5,8C5,10.21,6.79,12,9,12z M20,9V7 c0-0.55-0.45-1-1-1h0c-0.55,0-1,0.45-1,1v2h-2c-0.55,0-1,0.45-1,1v0c0,0.55,0.45,1,1,1h2v2c0,0.55,0.45,1,1,1h0c0.55,0,1-0.45,1-1 v-2h2c0.55,0,1-0.45,1-1v0c0-0.55-0.45-1-1-1H20z"
            />
          </g>
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
