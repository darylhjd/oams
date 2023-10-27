import axios from "axios";
import { LoginResponse } from "./login";
import { UserMeResponse } from "./user";

export class APIClient {
  static _client = axios.create({
    withCredentials: true,
    baseURL: `${process.env.API_SERVER}/api/v1`,
  });

  static readonly _loginPath = "/login";
  static readonly _userMePath = "/users/me";

  static async login(redirectUrl: string = ""): Promise<string> {
    const { data } = await this._client.get<LoginResponse>(this._loginPath, {
      params: {
        redirect_url: redirectUrl ? redirectUrl : process.env.WEB_SERVER,
      },
    });
    return data.redirect_url;
  }

  static async userMe(): Promise<UserMeResponse> {
    const { data } = await this._client.get<UserMeResponse>(this._userMePath);
    return data;
  }
}
