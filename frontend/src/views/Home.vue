<template>
  <div class="home">
    <!-- Nested objects are a bit weird rn, but we'll see -->
    <!-- <User
      v-bind="{
        bio: 'OH MY GOD I LOVE COCK OH GOD OH FUCK AAAAAAAAAAA AAAAAAA',
        createdAt: '2021-12-04T10:46:11.265Z',
        dweets: [
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'duderinos',
            id: 'sREiME3iw6',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:47.393Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:47.393Z',
            redweetCount: 1,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'my',
            id: 'SYqNUB1J0r',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:41.521Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:41.521Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'epic',
            id: '2eUnRZ9XJN',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:37.767Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:37.767Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'is',
            id: 'DrWVkyzWsJ',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:34.228Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:34.228Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'cock',
            id: 'yS1IKn5xbA',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:30.748Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:30.748Z',
            redweetCount: 0,
            replyCount: 0
          }
        ],
        email: '',
        feedObjects: [],
        followerCount: 1,
        followers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:47:38.95Z',
            email: 'okman@ok.com',
            followerCount: 0,
            followingCount: 1,
            name: 'OOOOK',
            pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayManALT'
          },
      
        ],
        following: [],
        followingCount: 0,
        likedDweets: [],
        name: 'OOOOK',
        pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
        redweetedDweets: [],
        redweets: [],
        username: 'OkayMan'
      }"
      viewUser="OkayManALT"
    />-->

    <!-- <div class="h-5"></div>

    <Dweet
      v-bind="{
        author: {
          bio: '',
          createdAt: '2021-12-04T10:46:11.265Z',
          email: 'ok@ok.com',
          followerCount: 1,
          followingCount: 0,
          name: 'OOOOK',
          pfpURL:
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
          username: 'OkayMan',
        },
        authorID: 'OkayMan',
        dweetBody: 'duderinos',
        id: 'sREiME3iw6',
        isReply: true,
        lastUpdatedAt: '2021-12-04T10:46:47.393Z',
        likeCount: 0,
        likeUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        media: [
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
        ],
        originalReplyID: '',
        postedAt: '2021-12-04T10:46:47.393Z',
        redweetCount: 1,
        redweetUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'RedweetManPagMan',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        replyCount: 0,
        replyDweets: [],
        replyTo: {
          author: {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
          authorID: 'OkayMan',
          dweetBody: 'my',
          id: 'sREiME3iw6',
          isReply: false,
          lastUpdatedAt: '2021-12-04T10:46:47.393Z',
          likeCount: 0,
          media: [
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
          ],
          originalReplyID: '',
          postedAt: '2021-12-04T10:46:47.393Z',
          redweetCount: 1,
          replyCount: 0,
        },
      }"
      viewUser="OkayMan"
    />

    <div class="h-5"></div>

    <FeedDweet
      v-bind="{
        author: {
          bio: '',
          createdAt: '2021-12-04T10:46:11.265Z',
          email: 'ok@ok.com',
          followerCount: 1,
          followingCount: 0,
          name: 'OOOOK',
          pfpURL:
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
          username: 'OkayMan',
        },
        authorID: 'OkayMan',
        dweetBody: 'duderinos',
        id: 'sREiME3iw6',
        isReply: true,
        lastUpdatedAt: '2021-12-04T10:46:47.393Z',
        likeCount: 0,
        likeUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        media: [
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
        ],
        originalReplyID: '',
        postedAt: '2021-12-04T10:46:47.393Z',
        redweetCount: 1,
        redweetUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'RedweetManPagMan',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        replyCount: 0,
        replyDweets: [],
        replyTo: {
          author: {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
          authorID: 'OkayMan',
          dweetBody: 'my',
          id: 'sREiME3iw6',
          isReply: false,
          lastUpdatedAt: '2021-12-04T10:46:47.393Z',
          likeCount: 0,
          media: [
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
          ],
          originalReplyID: '',
          postedAt: '2021-12-04T10:46:47.393Z',
          redweetCount: 1,
          replyCount: 0,
        },
      }"
      viewUser="OkayMan"
    />

    <div class="h-5"></div>

    <FeedRedweet
      v-bind="{
        author: {
          bio: '',
          createdAt: '2021-12-04T10:46:11.265Z',
          email: 'ok@ok.com',
          followerCount: 0,
          followingCount: 0,
          name: 'OOOOK',
          pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
          username: 'OkayMan'
        },
        authorID: 'OkayMan',
        originalRedweetID: 'sREiME3iw6',
        redweetOf: {
          author: {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 0,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan'
          },
          authorID: 'OkayMan',
          dweetBody: 'duderinos',
          id: 'sREiME3iw6',
          isReply: false,
          lastUpdatedAt: '2021-12-04T10:46:47.393Z',
          likeCount: 0,
          media: [
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
          ],
          originalReplyID: '',
          postedAt: '2021-12-04T10:46:47.393Z',
          redweetCount: 0,
          replyCount: 0
        },
        redweetTime: '2021-12-04T10:47:02.883Z'
      }"
      viewUser="OkayMan"
    />

    <div class="h-5"></div>

    <UserListObject
      v-bind="{
        bio: 'OH MY GOD I LOVE COCK OH GOD OH FUCK AAAAAAAAAAA AAAAAAA',
        createdAt: '2021-12-04T10:46:11.265Z',
        dweets: [
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'duderinos',
            id: 'sREiME3iw6',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:47.393Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:47.393Z',
            redweetCount: 1,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'my',
            id: 'SYqNUB1J0r',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:41.521Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:41.521Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'epic',
            id: '2eUnRZ9XJN',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:37.767Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:37.767Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'is',
            id: 'DrWVkyzWsJ',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:34.228Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:34.228Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'cock',
            id: 'yS1IKn5xbA',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:30.748Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:30.748Z',
            redweetCount: 0,
            replyCount: 0
          }
        ],
        email: '',
        feedObjects: [],
        followerCount: 1,
        followers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:47:38.95Z',
            email: 'okman@ok.com',
            followerCount: 0,
            followingCount: 1,
            name: 'OOOOK',
            pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayManALT'
          },
      
        ],
        following: [],
        followingCount: 0,
        likedDweets: [],
        name: 'OOOOK',
        pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
        redweetedDweets: [],
        redweets: [],
        username: 'OkayMan'
      }"
      viewUser="OkayManALT"
    />

    <div class="h-5"></div>

    <EditDweet
      v-bind="{
        author: {
          bio: '',
          createdAt: '2021-12-04T10:46:11.265Z',
          email: 'ok@ok.com',
          followerCount: 1,
          followingCount: 0,
          name: 'OOOOK',
          pfpURL:
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
          username: 'OkayMan',
        },
        authorID: 'OkayMan',
        dweetBody: 'duderinos',
        id: 'sREiME3iw6',
        isReply: true,
        lastUpdatedAt: '2021-12-04T10:46:47.393Z',
        likeCount: 0,
        likeUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        media: [
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
          'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
        ],
        originalReplyID: '',
        postedAt: '2021-12-04T10:46:47.393Z',
        redweetCount: 1,
        redweetUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'RedweetManPagMan',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        replyCount: 0,
        replyDweets: [],
        replyTo: {
          author: {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
          authorID: 'OkayMan',
          dweetBody: 'my',
          id: 'sREiME3iw6',
          isReply: false,
          lastUpdatedAt: '2021-12-04T10:46:47.393Z',
          likeCount: 0,
          media: [
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FgO25ZZgYrftxcwbNpydD5AdDkjijSW.mp4?generation=1640258723535975&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2FBmq2bpVF50I25GeJcZ8cLdio4KTnn8.png?generation=1640258942272122&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fl4kV5WBkWsXI7AbRlkxjj0uiuHNZ3G.gif?generation=1640247339114616&alt=media',
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/media%2Fefy7lIUbTUwCPNUjICnaxsmpdYqmb1.jpg?generation=1640259076665041&alt=media'
          ],
          originalReplyID: '',
          postedAt: '2021-12-04T10:46:47.393Z',
          redweetCount: 1,
          replyCount: 0,
        },
      }"
      viewUser="OkayMan"
    />-->
    <EditUser
      v-bind="{
        bio: 'OH MY GOD I LOVE COCK OH GOD OH FUCK AAAAAAAAAAA AAAAAAA',
        createdAt: '2021-12-04T10:46:11.265Z',
        dweets: [],
        email: 'epicman@epic.com',
        feedObjects: [],
        followerCount: 1,
        followers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:47:38.95Z',
            email: 'okman@ok.com',
            followerCount: 0,
            followingCount: 1,
            name: 'OOOOK',
            pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayManALT'
          },
        ],
        following: [],
        followingCount: 0,
        likedDweets: [],
        name: 'OOOOK',
        pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
        redweetedDweets: [],
        redweets: [],
        username: 'OkayMan'
      }"
      viewUser="OkayManALT"
    />
  </div>
</template>

<script>
// @ is an alias to /src
import User from "../components/User.vue";
import Dweet from "../components/Dweet.vue";
import FeedDweet from "../components/FeedDweet.vue";
import FeedRedweet from "../components/FeedRedweet.vue";
import UserListObject from "../components/UserListObject.vue";
import EditDweet from "../components/EditDweet.vue";
import EditUser from "../components/EditUser.vue";

export default {
  name: "Home",
  components: {
    User,
    Dweet,
    FeedDweet,
    FeedRedweet,
    UserListObject,
    EditDweet,
    EditUser,
  },
};
</script>
