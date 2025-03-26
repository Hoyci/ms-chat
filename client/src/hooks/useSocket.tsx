import { useEffect } from "react";
import SocketService from "@services/socket";
import { config } from "config";

export function useSocket() {
  const token =
    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IiIsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiam9obmRvZSIsImVtYWlsIjoiam9obmRvZUBlbWFpbC5jb20iLCJpc3MiOiJtcy1jaGF0LWF1dGgiLCJleHAiOjE3NDI5OTM3NjIsImlhdCI6MTc0Mjk5MDE2MiwianRpIjoiY2I4ZDhmMmMtZmYxNC00ZTJhLThkOWUtMGZhNGM5YzFmYzEzIn0.LD8Vxyos0gc2ahV-jYshZnFtBituRaEnJgZun1VXS820ZfJcSSQz41-oUySbhk_2MSxpmp1I3ypAu1hZWFgOaBCYB5jyYApaAllSnJM2rw5VKuwpkRjzZpmBzgMVQvHMq8V-iaDs3LCwkNDJfPYYj_WQliklikhN3EZs2-24AKq8iYqWm75VZyCG47eFhDjtd7gR3zDhk-YDkc0iB4E4EBa2SnyQm0N3G5jbcCbcda11UvvRm0LjPmxUeNc8WQ6lrZSaWgelSzkvCAPr17zxfQ2j9NbA4F4Sw0oa0TToevutB2lufSYX5r5lXUaoUnj9Qq6Jw8kmEW3E8EFrhJZlbA";
  useEffect(() => {
    SocketService.connect(`${config.API_URL}/api/v1/ws?token=${token}`);

    return () => {
      SocketService.disconnect();
    };
  }, []);
}
