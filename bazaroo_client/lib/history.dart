import 'package:flutter/material.dart';
import 'nav/customer_nav.dart'; 
import 'homecustomer.dart';

class PurchaseHistory extends StatelessWidget {
  final String userId;
  PurchaseHistory({required this.userId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80, 
        backgroundColor: Colors.red,
        title: Row(
          children: [
            Text(
              'Purchase History',
              style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: Icon(Icons.chevron_left, color: Colors.white, size: 35),
              onPressed: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => HomeScreen(userId: userId)), 
                );
              },
            ),
          ),
        ],
      ),
      body: Center(
        child: Text('This is the Purchase History screen'),
      ),
      bottomNavigationBar: BottomNavBar(userId: userId),
    );
  }
}
