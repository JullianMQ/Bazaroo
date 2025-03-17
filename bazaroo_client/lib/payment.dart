import 'package:flutter/material.dart';
import 'homecustomer.dart';

class PaymentOptions extends StatelessWidget {
  final String userId;
  PaymentOptions({required this.userId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        toolbarHeight: 80, 
        title: Row(
          children: [
            Icon(Icons.credit_card, color: Colors.red, size: 35), 
            SizedBox(width: 8),
            Text(
              'Payment Options',
              style: TextStyle(color: Colors.red, fontSize: 24, fontWeight: FontWeight.bold),
            ),
          ],
        ),
        automaticallyImplyLeading: false,
        actions: [
          Container(
            margin: EdgeInsets.symmetric(horizontal: 25.0),
            child: IconButton(
              icon: Icon(Icons.chevron_left, color: Colors.red, size: 35),
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
        child: Text('This is the Payment screen'),
      ),
    );
  }
}
