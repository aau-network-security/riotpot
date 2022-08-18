import { InteractionOption, NetworkOption } from "../../constants/globals";

export interface Service {
  name: string;
  interaction: InteractionOption;
  network: NetworkOption;
  host: string;
  port: Number;
}
