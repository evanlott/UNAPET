#include <iostream>
#include <cmath>
#include <vector>
#include <thread>
#include <future>
#include <chrono>	
#include <mutex>
using namespace std;

int main()
{
	int x = 0;
	this_thread::sleep_for(chrono::seconds(0));
	cout << (5/x);
}
	