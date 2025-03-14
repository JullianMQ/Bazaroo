import 'package:flutter/material.dart';
import 'cart.dart';
import 'categories.dart';
import 'address.dart';
import 'payment.dart';
import 'history.dart';
import 'nav/bottom_nav.dart';  
import 'main.dart';

class HomeScreen extends StatefulWidget {
  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
    @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        toolbarHeight: 80,
        actions: <Widget>[Container()],
        iconTheme: IconThemeData(color: Colors.white),
        backgroundColor: Colors.red,
        title: Row(
          children: [
           GestureDetector(
              onTap: () {
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => MyApp()),  
                );
              },
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 8.0),
                child: Row(
                  children: [
                    Image.asset(
                      'assets/Bazaroo-logo-white.png',
                      height: 40,
                    ),
                    Image.asset(
                      'assets/Bazaroo-text-white.png',
                      height: 30,
                    ),
                  ],
                ),
              ),
            ),
            Spacer(),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 8.0),
              child: Builder(
                builder: (context) => IconButton(
                  icon: Icon(Icons.shopping_cart, color: Colors.white, size: 25),
                  onPressed: () => Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => Cart()),
                  ),
                ),
              ),
            ),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 8.0),
              child: Builder(
                builder: (context) => IconButton(
                  icon: Icon(Icons.menu, color: Colors.white, size: 25),
                  onPressed: () => Scaffold.of(context).openEndDrawer(),
                ),
              ),
            ),
          ],
        ),
      ),
      endDrawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            SizedBox(
              height: 90,
              child: DrawerHeader(
                child: Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: EdgeInsets.symmetric(horizontal: 10.0, vertical: 5.0),
                    child: IconButton(
                      icon: Icon(Icons.close, color: Colors.red),
                      onPressed: () => Navigator.of(context).pop(),
                    ),
                  ),
                ),
              ),
            ),
            ListTile(title: Text('My Account')),
            ListTile(title: Text('Purchase History'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => PurchaseHistory()));
            }),
            ListTile(title: Text('Addresses'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => AddressScreen()));
            }),
            ListTile(title: Text('Payment Options'), onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (context) => PaymentOptions()));
            }),
            ListTile(title: Text('Change Password')),
            ListTile(title: Text('Settings')),
            ListTile(title: Text('Log out')),
          ],
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: EdgeInsets.all(20.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Featured
              Container(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container( // Heading
                      margin: EdgeInsets.only(bottom: 10),
                      child: Text('Featured',
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          color: Colors.red,
                        ),
                      ),
                    ),
                    Container( // Contents or image placeholder
                      margin: EdgeInsets.only(bottom: 40),
                      child: AspectRatio(
                        aspectRatio: 16 / 9,
                        child: Container(
                          decoration: BoxDecoration(
                            color: Colors.grey[300],
                            borderRadius: BorderRadius.circular(15),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              // Categories
              Container(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container( // Heading
                      margin: EdgeInsets.only(bottom: 10),
                      child: GestureDetector(
                        onTap: () {
                          // Navigate to Categories screen
                          Navigator.push(
                            context,
                            MaterialPageRoute(builder: (context) => Categories()),
                          );
                        },
                        child: Text(
                          'Categories',
                          style: TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                            color: Colors.red,
                          ),
                        ),
                      ),
                    ),
                    Container( // Contents
                      child: GridView.builder(
                        shrinkWrap: true,
                        physics: NeverScrollableScrollPhysics(),
                        gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2,
                          crossAxisSpacing: 16,
                          mainAxisSpacing: 16,
                          childAspectRatio: 1,
                        ),
                        itemCount: 8,
                        itemBuilder: (context, index) {
                          return Container(
                            decoration: BoxDecoration(
                              color: Colors.grey[300],
                              borderRadius: BorderRadius.circular(15),
                            ),
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: BottomNavBar(
      ),
    );
  }
}
