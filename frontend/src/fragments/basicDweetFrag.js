import gql from "graphql-tag";
import { basicUserFrag } from "./basicUserFrag";

export const basicDweetFrag = gql`
  fragment BasicDweetFrag on BasicDweet {
    dweetBody
    id
    author {
      ...BasicUserFrag
    }
    authorID
    postedAt
    lastUpdatedAt
    likeCount
    isReply
    originalReplyID
    replyCount
    redweetCount
    media
  }
  ${basicUserFrag}
`
