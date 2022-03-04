import gql from "graphql-tag";
import { basicDweetFrag } from "./basicDweetFrag";
import { redweetFrag } from "./redweetFrag";

export const basicFeedObjectFrag = gql`
  fragment BasicFeedObjectFrag on BasicFeedObject {
    ...BasicDweetFrag
    ...RedweetFrag
  }
  ${basicDweetFrag}
  ${redweetFrag}
`
