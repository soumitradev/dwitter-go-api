import gql from "graphql-tag";

export const basicUserFrag = gql`
  fragment BasicUserFrag on BasicUser {
    username
    name
    email
    bio
    pfpURL
    followerCount
    followingCount
    createdAt
  }
`
