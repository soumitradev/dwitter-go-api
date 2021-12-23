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
      <div
        class="flex flex-row self-start p-4"
        v-if="followers.map(e => e.username).includes(viewUser)"
      >
        <button
          class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-secondary-90 hover:text-secondary-40 transition duration-200 ease-in-out"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 fill-current" viewBox="0 0 24 24">
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
      <div class="flex flex-row self-start p-4" v-else>
        <button
          class="p-2 rounded-full text-neutralVariant-50 bg-neutral-99 bg-opacity-20 hover:bg-opacity-30 hover:bg-primary-90 hover:text-primary-40 transition duration-200 ease-in-out"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 fill-current" viewBox="0 0 24 24">
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

    <span class="text-left px-2 pt-6 text-neutral-10">{{ bio }}</span>

    <div class="flex flex-row mt-2 text-neutralVariant-50" v-if="email">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5 mx-1 fill-current"
        viewBox="0 0 24 24"
      >
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path
          d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-.4 4.25l-7.07 4.42c-.32.2-.74.2-1.06 0L4.4 8.25c-.25-.16-.4-.43-.4-.72 0-.67.73-1.07 1.3-.72L12 11l6.7-4.19c.57-.35 1.3.05 1.3.72 0 .29-.15.56-.4.72z"
        />
      </svg>
      <span>{{ email }}</span>
    </div>
    <div class="flex flex-row mt-2 text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5 mx-1 fill-current"
        viewBox="0 0 24 24"
      >
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path
          d="M20 3h-1V2c0-.55-.45-1-1-1s-1 .45-1 1v1H7V2c0-.55-.45-1-1-1s-1 .45-1 1v1H4c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 18H5c-.55 0-1-.45-1-1V8h16v12c0 .55-.45 1-1 1z"
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
        class="h-4 w-4 mr-2 fill-current"
        viewBox="0 0 24 24"
      >
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path
          d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V18c0 .55.45 1 1 1h12c.55 0 1-.45 1-1v-1.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05.02.01.03.03.04.04 1.14.83 1.93 1.94 1.93 3.41V18c0 .35-.07.69-.18 1H22c.55 0 1-.45 1-1v-1.5c0-2.33-4.67-3.5-7-3.5z"
        />
      </svg>
      <span v-if="following.length > 0">Follows {{ overflowNames(following) }}</span>
      <span v-else>Doesn't follow anyone you follow</span>
    </div>
    <div class="flex flex-row mb-2 text-left text-sm text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-4 w-4 mr-2 fill-current"
        viewBox="0 0 24 24"
      >
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path
          d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V18c0 .55.45 1 1 1h12c.55 0 1-.45 1-1v-1.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05.02.01.03.03.04.04 1.14.83 1.93 1.94 1.93 3.41V18c0 .35-.07.69-.18 1H22c.55 0 1-.45 1-1v-1.5c0-2.33-4.67-3.5-7-3.5z"
        />
      </svg>
      <span v-if="followers.length > 0">Followed by {{ overflowNames(followers) }}</span>
      <span v-else>Isn't followed by anyone you follow</span>
    </div>
  </div>
</template>

<script>

export default {
  name: "User",
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
