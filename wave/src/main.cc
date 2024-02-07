#include <iostream>
#include <boost/asio.hpp>


#include "shared/logger/logger.hpp"

using namespace std;

int main() {
	boost::asio::io_context io_context;
	logger::LOGLEVEL;
	cout << "Hallo" << endl;
	return 0;
}
