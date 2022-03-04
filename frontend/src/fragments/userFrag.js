import gql from "graphql-tag";
import { basicUserFrag } from "./basicUserFrag";
import { basicDweetFrag } from "./basicDweetFrag";
import { redweetFrag } from "./redweetFrag";
import { basicFeedObjectFrag } from "./basicFeedObjectFrag";

export const userFrag = gql`
  fragment UserFrag on User {
    username
    name
    email
    bio
    pfpURL
    dweets {
      ...BasicDweetFrag
    }
    redweets {
      ...RedweetFrag
    }
    redweetedDweets {
      ...BasicDweetFrag
    }
    feedObjects {
      ...BasicFeedObjectFrag
    }
    likedDweets {
      ...BasicDweetFrag
    }
    followerCount
    followers {
      ...BasicUserFrag
    }
    followingCount
    following {
      ...BasicUserFrag
    }
    createdAt
  }
  ${basicDweetFrag}
  ${redweetFrag}
  ${basicFeedObjectFrag}
  ${basicUserFrag}
`
