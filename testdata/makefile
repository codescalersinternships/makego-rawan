hello: hello1.txt hello2.txt 
	cat hello1.txt
	cat hello2.txt
	
hello1.txt: 
	echo "Hello 1, Make!" > hello1.txt

hello2.txt: 
	echo "Hello 2, Make!" > hello2.txt

clean:
	rm -f hello1.txt
	rm -f hello2.txt
