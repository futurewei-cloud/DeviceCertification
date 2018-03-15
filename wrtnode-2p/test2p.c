#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <errno.h>
#include <malloc.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/ioctl.h>
#include <stdarg.h>
#include <fcntl.h>
#include <fcntl.h>
#include <termios.h>

#define BUF_LEN 512

int main(int argc, char **argv){
	struct termios serialio;
	char * write_buf;
	unsigned char buf[BUF_LEN] = {0};
	const char device [] = "/dev/ttyS1";
	int len, fd, i;
	fd = open(device, O_RDWR | O_NOCTTY);
	if(fd < 0){
		printf("open %s failed\n", device);
		return -1;
	}

	/* setup tty */
	if(tcgetattr(fd, &serialio) < 0)
		printf("warning: get %s information failed\n", device);
	/* set baud to 115200 */
	cfsetispeed(&serialio, B115200);
	cfsetospeed(&serialio, B115200);

	serialio.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);
	cfmakeraw(&serialio);
	
	if(tcsetattr(fd, TCSANOW, &serialio) < 0)
		printf("warning: set %s information failed\n", device);

#if 1
	write_buf = malloc(argc);

        len = 0;
	for(i=0;i<(argc - 1);i++,len++){
		write_buf[i] = strtol(argv[i+1], NULL, 16);
	}

	for(i=0; i<len; i++)
		printf(" %02x", write_buf[i]);

	printf("\n");

	if(argc < 6)
		return -1;

	len = write(fd, write_buf, len);
	printf("write %d\n", len);
#endif

	usleep(100 * 1000);

	len = read(fd, buf, BUF_LEN);
	printf("read %d\n", len);

	printf("buf:");
	for(i=0; i<len; i++)
		printf(" %02x", buf[i]);

	printf("\n");
	return 0;

}
