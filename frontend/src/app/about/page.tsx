import { Space, Text } from "@mantine/core";
import styles from "@/styles/AboutPage.module.css";

export default function AboutPage() {
  return (
    <>
      <Text ta="center" className={styles.aboutText}>
        The Online Attendance Management System provides a centralised
        attendance service for Nanyang Technological University.
      </Text>
      <Space h="xl" />
      <Text ta="center" size="sm">
        Copyright &#169; Har Jing Daryl, 2023
      </Text>
    </>
  );
}
