import 'package:flutter/material.dart';
import 'logincustomer.dart';
import 'loginseller.dart';

class StartupScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: Column(
        children: [

          SizedBox(height: 80),

          Center(
            child: Image.asset(
              'assets/Bazaroo-full-logo-red.png',
              height: 350, 
            ),
          ),
          SizedBox(height: 30),
          
          Text(
            'Are you a...',
            style: TextStyle(
              fontSize: 24,
              color: Colors.black,
            ),
          ),
          SizedBox(height: 40),
          
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => LoginCustomer()),
                  );
                },
                child: Column(
                  children: [
                    Icon(
                      Icons.person,
                      size: 100,
                      color: Colors.red,
                    ),
                    Text(
                      'Customer',
                      style: TextStyle(
                        fontSize: 18,
                        color: Colors.red,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(width: 50),

              GestureDetector(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => LoginSeller()),
                  );
                },
                child: Column(
                  children: [
                    Icon(
                      Icons.work,
                      size: 100,
                      color: Colors.red,
                    ),
                    Text(
                      'Seller',
                      style: TextStyle(
                        fontSize: 18,
                        color: Colors.red,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}