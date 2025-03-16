import 'package:bazaroo_client/startupscreen.dart';
import 'package:flutter/material.dart';
import 'nav/seller_nav.dart';

class SellerHome extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80, 
        backgroundColor: Colors.red,
        title: Row(
          children: [
            Text(
              'Seller home page',
              style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          GestureDetector(
            onTap: () {
              Navigator.pushReplacement(
                context,
                MaterialPageRoute(builder: (context) => StartupScreen()), // Redirect to StartupScreen after logout
              );
            },
            child: Row(
              children: [
                Padding(
  padding: EdgeInsets.symmetric(horizontal: 20.0), // Adjust padding as needed
  child: GestureDetector(
    onTap: () {
      // Add your logout functionality here
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => StartupScreen()), // Redirect to StartupScreen after logout
      );
    },
    child: Row(
      children: [
        Icon(Icons.logout, color: Colors.white, size: 30),
        SizedBox(width: 8),
        Text(
          'Logout',
          style: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold),
        ),
      ],
    ),
  ),
)

              ],
            ),
          )

        ],
      ),
      body: Center(
        child: Text('This is the sellers home page'),
      ),
      bottomNavigationBar: BottomNavBar(
      ),
    );
  }
}