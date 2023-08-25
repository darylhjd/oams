import { MediaQuery } from "@mantine/core";

// Shows children in mobile mode.
export function Mobile({ children }: { children: React.ReactNode }) {
  return (
    <MediaQuery largerThan='md' styles={{ display: 'none' }}>
      {children}
    </MediaQuery>
  )
}

// Shows children in desktop mode.
export function Desktop({ children }: { children: React.ReactNode }) {
  return (
    <MediaQuery smallerThan='md' styles={{ display: 'none' }}>
      {children}
    </MediaQuery>
  )
}
