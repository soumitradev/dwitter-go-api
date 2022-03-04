import gql from "graphql-tag";
import { basicUserFrag } from "./basicUserFrag";
import { basicDweetFrag } from "./basicDweetFrag";

export const dweetFrag = gql`
  fragment DweetFrag on Dweet {
    dweetBody
    id
    author {
      ...BasicUserFrag
    }
    authorID
    postedAt
    lastUpdatedAt
    likeCount
    likeUsers {
      ...BasicUserFrag
    }
    isReply
    originalReplyID
    replyTo {
      ...BasicDweetFrag
    }
    replyCount
    replyDweets {
      ...BasicDweetFrag
    }
    redweetCount
    redweetUsers {
      ...BasicUserFrag
    }
    media
  }
  ${basicUserFrag}
  ${basicDweetFrag}
`
