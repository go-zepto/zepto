import React from "react";
import { Layout } from "react-admin";
import { Schema } from "../../../types/schema";
import MenuGenerator from "./compGens/Menu/menu";


export const generateLayoutComp = (s: Schema): React.FC => {
  const Menu = MenuGenerator(s);
  return (props: any) => (
    <Layout {...props} menu={Menu}  />
  );
};
