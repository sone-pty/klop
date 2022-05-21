#include <pthread.h>
#include <stdio.h>
#include <unistd.h>
#include <strings.h>
#include <string.h>
#include <vector>

#include <ds/uf.h>

#include <http_parser/http_parser.h>

const char *HTTPREQ_1 = "GET /search?hl=zh-CN&source=hp&q=domety&aq=f HTTP/1.1\r\nHost: 127.0.0.1\r\nContent-Length: 5\r\n\r\nHello";
const char *HTTPREQ_2 = "POST /upload?fileName=1.txt HTTP/1.1\r\nHost: 127.0.0.1\r\nContent-Length: 5\r\n\r\nHello";

int onMessageComplete(http_parser *parser)
{
    if (HTTP_PARSER_ERRNO(parser) == HPE_OK)
    {
        printf("Message parsed successfully\n");
    }

    return 0;
}

int main(void)
{
    /*
    http_parser parser;
    http_parser_settings parser_settings;

    http_parser_init(&parser, HTTP_REQUEST);
    http_parser_settings_init(&parser_settings);

    parser_settings.on_message_complete = onMessageComplete;

    http_parser_execute(&parser, &parser_settings, HTTPREQ_1, strlen(HTTPREQ_1));
    http_parser_execute(&parser, &parser_settings, HTTPREQ_2, strlen(HTTPREQ_2));
    */

    return 0;
}