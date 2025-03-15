import 'package:flutter/material.dart';
import 'nav/bottom_nav.dart';

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
        child: Text('This is the sellers home page'),
      ),
      bottomNavigationBar: BottomNavBar(
      ),
    );
  }
}
