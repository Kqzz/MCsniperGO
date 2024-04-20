import { extendTheme } from "@chakra-ui/react";

const config = {
  initialColorMode: "dark",
  useSystemColorMode: false,
};

export const chakraCustomTheme = extendTheme({
  config,
});

console.log(config);
