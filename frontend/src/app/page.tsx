import { APIClient } from "@/api/client"

export default async function Page() {
  const user = await APIClient.getUserMe()
  console.log(user)
  console.log("Building")

  return (
    <div>
      This is the home page.
    </div>
  )
}
