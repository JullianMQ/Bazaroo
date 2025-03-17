import 'package:flutter/material.dart';
import 'nav/customer_nav.dart';

class Cart extends StatelessWidget {
  final String userId;
  Cart({required this.userId});
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80, 
        backgroundColor: Colors.red,
        title: Row(
          children: [
            Text(
              'Cart',
              style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0, ),
            child: IconButton(
              icon: Icon(Icons.chevron_left, color: Colors.white, size: 35,),
              onPressed: () {
                Navigator.pop(context);
              },
            ),
          ),
        ],
      ),
      body: Center(
        child: Text('This is the Cart screen'),
      ),
      bottomNavigationBar: BottomNavBar(userId: userId),
    );
  }
}
