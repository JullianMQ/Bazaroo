import 'package:bazaroo_client/index.dart';
import 'package:flutter/material.dart';
import 'nav/seller_nav.dart';

class SellerHome extends StatelessWidget {
  final String userId;
  const SellerHome({Key? key, required this.userId}) : super(key: key);

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
                MaterialPageRoute(builder: (context) => StartupScreen()),
              );
            },
            child: Row(
              children: [
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: 20.0), 
                  child: GestureDetector(
                    onTap: () {
                      Navigator.pushReplacement(
                        context,
                        MaterialPageRoute(builder: (context) => StartupScreen()),
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
        child: Text('This is the sellers home page $userId'),
      ),
      bottomNavigationBar: BottomNavBar(userId: userId),
    );
  }
}