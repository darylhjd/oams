import axios from "axios"
import { User } from "./models"

export class APIClient {
  static _client = axios.create({
    baseURL: `${process.env.API_SERVER_HOST}:${process.env.API_SERVER_PORT}/api/v1`,
    withCredentials: true,
  })

  static readonly _userMePath = "/users/me"

  static async getUserMe(): Promise<User | null> {
    try {
      const { data } = await this._client.get<User>(this._userMePath)
      return data
    } catch (error) {
      return null
    }
  }
}
