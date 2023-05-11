import {RootRoute, Route} from '@tanstack/router';
import './index.css';
import {Box, Center, Flex, Text} from '@chakra-ui/react';

const path = "/"

export function Component() {
  return (
    <Center>
      <Box>
        <Flex>
          <Text>OAMS: Online Attendance Management System</Text>
        </Flex>
      </Box>
    </Center>
  )
}

function IndexPage(root: RootRoute) {
  return new Route({
    getParentRoute: () => root,
    path: path,
    component: Component,
  })
}

export default IndexPage
