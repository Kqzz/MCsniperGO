import React, { useState } from "react";
import {
  Box,
  Checkbox,
  FormControl,
  FormLabel,
  Input,
  Button,
  HStack,
} from "@chakra-ui/react";

function ClaimForm({ queueClaim }) {
  const [username, setUsername] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [useProxies, setUseProxies] = useState(true);
  const [runInfinitely, setRunInfinitely] = useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();
    if (runInfinitely) {
      setEndTime(-1);
      setStartTime(-1);
    }

    queueClaim({
      username,
      startTime: new Number(startTime),
      endTime: new Number(endTime),
      useProxies,
    });
  };

  const handleStartTimeChange = (e) => {
    setStartTime(new Number(e.target.value));
  };

  const handleEndTimeChange = (e) => {
    setEndTime(new Number(e.target.value));
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString(); // Adjust the format as per your requirements
  };

  return (
    <Box width="100%" borderWidth="1px" borderRadius="lg" p={4}>
      <form onSubmit={handleSubmit}>
        <HStack spacing={4} align="flex-end">
          <FormControl isRequired>
            <FormLabel>Username</FormLabel>
            <Input
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </FormControl>

          <FormControl isRequired isDisabled={runInfinitely}>
            <FormLabel>Start Time (Unix Timestamp)</FormLabel>
            <Input
              type="text"
              value={startTime}
              onChange={handleStartTimeChange}
            />
          </FormControl>

          <FormControl isRequired isDisabled={runInfinitely}>
            <FormLabel>End Time (Unix Timestamp)</FormLabel>
            <Input type="text" value={endTime} onChange={handleEndTimeChange} />
          </FormControl>
        </HStack>
        <HStack mt={3}>
          <FormControl>
            <FormLabel>Start Time</FormLabel>
            <Input type="text" value={formatTimestamp(startTime)} isReadOnly />
          </FormControl>

          <FormControl>
            <FormLabel>End Time</FormLabel>
            <Input type="text" value={formatTimestamp(endTime)} isReadOnly />
          </FormControl>
        </HStack>
        <HStack mt={3}>
          <Checkbox
            isChecked={useProxies}
            onChange={(e) => setUseProxies(e.target.checked)}
          >
            Use proxies
          </Checkbox>

          <Checkbox
            isChecked={runInfinitely}
            onChange={(e) => setRunInfinitely(e.target.checked)}
          >
            Run infinitely
          </Checkbox>
          <Button type="submit" colorScheme="teal">
            Queue Claim
          </Button>
        </HStack>
      </form>
    </Box>
  );
}

export default ClaimForm;
