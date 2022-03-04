import gql from "graphql-tag";
import { basicUserFrag } from "./basicUserFrag";
import { basicDweetFrag } from "./basicDweetFrag"

export const redweetFrag = gql`
  fragment RedweetFrag on Redweet {
    author {
      ...BasicUserFrag
    }
    authorID
    redweetOf {
      ...BasicDweetFrag
    }
    originalRedweetID
    redweetTime
  }
  ${basicUserFrag}
  ${basicDweetFrag}
`
