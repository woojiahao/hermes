import {JsonConverter, JsonCustomConvert} from "json2typescript"
import {formatDate} from "../utility/general"

@JsonConverter
export class DateConverter implements JsonCustomConvert<Date> {
  serialize(data: Date): any {
    return formatDate(data)
  }

  deserialize(data: any): Date {
    return new Date(data)
  }
}
