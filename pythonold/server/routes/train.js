const express = require('express');
const path = require('path');
const fs = require('fs-extra');
let router = express.Router();
const testFolder = './public/archives/';
const faceFolder = './public/faces/';
router.post('/faces', function (req, res) {
    //console.log(req.body)
    Object.keys(req.body).forEach(key => {
        if (key === "inconnu")
            return
        req.body[key].forEach(img => {
            var oldFile = undefined;
            var newFile = "";
            if (img.origine === undefined) {
                oldFile = testFolder + img.name
            }
            if (img.origine !== undefined && img.origine !== key) {
                oldFile = faceFolder + img.origine + '/' + img.name
            }
            newFile = faceFolder + key + '/' + img.name

            if (oldFile !== undefined) {
                fs.ensureDir(faceFolder + key)
                    .then(() => {
                        fs.rename(oldFile, newFile, function (err) {
                            if (err) {
                                console.log(err)
                            }
                        });
                    })
                    .catch(err => {
                        console.log(err)
                    })
            }
        })

    })
    res.status(200).send({})
})
router.get('/faces', function (req, res) {
    result = {}

    var walk = function (dir, results, done) {
        fs.readdir(dir, function (err, list) {
            if (err) return done(err);
            var i = 0;
            (function next() {
                var item = list[i++];
                if (!item) return done(null, results);

                file = dir + '/' + item;
                if (dir !== faceFolder && !results[dir.replace(faceFolder, "").replace('/', '')]) results[dir.replace(faceFolder, "").replace('/', '')] = []
                fs.stat(file, function (err, stat) {
                    if (stat && stat.isDirectory()) {
                        walk(file, results, function (err, res) {
                            // results = results.concat(res);
                            next();
                        });
                    } else {
                        results[dir.replace(faceFolder, "").replace('/', '')].push({ name: item, origine: dir.replace(faceFolder, "").replace('/', '') });
                        next();
                    }
                });
            })();
        });
    };

    walk(faceFolder, {}, function (err, result) {
        fs.readdir(testFolder, (err, files) => {
            array = []
            files.forEach(file => {
                array.push({ name: file });
            });
            result['inconnu'] = array

            console.log(err, result)

            res.status(200).send(result)
        })

    })


    //      res.status(200).send(array)

})
module.exports = router;