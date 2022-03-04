import gql from "graphql-tag";
import { dweetFrag } from "./dweetFrag";
import { redweetFrag } from "./redweetFrag";

export const feedObjectFrag = gql`
  fragment FeedObjectFrag on FeedObject {
    ...DweetFrag
    ...RedweetFrag
  }
  ${dweetFrag}
  ${redweetFrag}
`
