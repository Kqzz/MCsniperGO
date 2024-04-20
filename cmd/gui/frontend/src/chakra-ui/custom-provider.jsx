import React from "react";
import { ChakraProvider } from "@chakra-ui/react";
import { chakraCustomTheme } from "./custom-theme";

export const ChakraUIProvider = ({ children }) => {
  return <ChakraProvider theme={chakraCustomTheme}>{children}</ChakraProvider>;
};
