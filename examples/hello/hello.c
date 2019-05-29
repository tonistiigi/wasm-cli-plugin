#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <dirent.h>
#include <unistd.h>
#include <fcntl.h>

int main() {
	printf("Hello world, I am %s!\n", getenv("WHOAMI"));
	
	DIR *dir;
	struct dirent *entry;
	
	if ((dir = opendir("/")) == NULL)
		perror("opendir() error");
	else {
		puts("contents of / :");
		while ((entry = readdir(dir)) != NULL)
			printf("  %d\t%s\n", entry->d_type, entry->d_name);
		closedir(dir);
	}
	
	printf("\nwriting to /foo \n");
	int out = open("/foo", O_WRONLY | O_CREAT, 0660);
	if (out >= 0)
	{
		const char* msg = "a message\n";
		write(out, msg, strlen(msg));
		close(out);
	} else {
		perror("failed to open /foo");
	}
	
	if ((dir = opendir("/")) == NULL)
		perror("opendir() error");
	else {
		puts("contents of / after write :");
		while ((entry = readdir(dir)) != NULL)
			printf("  %d\t%s\n", entry->d_type, entry->d_name);
		closedir(dir);
	}
	
	return 0;
}