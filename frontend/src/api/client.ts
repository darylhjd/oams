import axios from "axios"
import { LoginResponse, UserMeResponse } from "./models"

export class APIClient {
  static _client = axios.create({
    baseURL: `${process.env.API_SERVER_HOST}:${process.env.API_SERVER_PORT}/api/v1`,
    withCredentials: true,
  })

  static readonly _loginPath = "/login"
  static readonly _logoutPath = "/logout"

  static readonly _userMePath = "/users/me"

  static async getLoginUrl(returnTo: string): Promise<string> {
    if (returnTo.length == 0) {
      returnTo = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}`
    }

    try {
      const { data } = await this._client.get<LoginResponse>(this._loginPath, {
        params: {
          redirect_url: returnTo,
        }
      })
      return data.redirect_url
    } catch (error) {
      throw error
    }
  }

  static async logout(): Promise<boolean> {
    try {
      await this._client.get(this._logoutPath)
      return true
    } catch (error) {
      return false
    }
  }

  static async getUserMe(): Promise<UserMeResponse | null> {
    try {
      const { data } = await this._client.get<UserMeResponse>(this._userMePath)
      return data
    } catch (error) {
      return null
    }
  }
}
