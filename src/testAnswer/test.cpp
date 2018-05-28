#include <stdio.h>
#include <string.h>

int main(){
	char str1[]="ACRFDJINM";
	char str2[] = "";
	int i;

	for(i=0;i<(strlen(str1));i++)
	{
		if((str1[i] > 0x40)&&(str1[i] < 0x5B))
			str2[i] = str1[i]+0x20;
		else
			str2[i] = str1[i];
	}
	str2[i]='\0';
	printf("\n -------%s--------\n",str1);
	printf("\n -------%s--------\n",str2);

	return 0;
}
